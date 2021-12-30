const { ShardingManager } = require('discord.js-light');
const config = require('./config');
const { printLog, getFiles } = require('./utils');
const { botCache } = require('./structures/cache');

const fetch = (...args) => import('node-fetch').then(({ default: fetch }) => fetch(...args));

const manager = new ShardingManager('app.js', { token: config.botToken });
manager.on('shardCreate', shard => printLog('Launched shard.', 'INFO', shard.id));

const main = async function () {
    // ========== GET COMMANDS ==========
    const cmdFiles = await getFiles('./commands');
    const commands = [];

    for (const file of cmdFiles) {
        const cmdFile = require(`./commands/${file}`);
        const cmdName = file.split('.')[0];

        // Set the command
        botCache.commands.set(cmdName, cmdFile.command);
        // Check if there are any buttons
        if (cmdFile.buttons != null) {
            // Loop over the buttons
            for (let i = 0; i < cmdFile.buttons.length; i++) {
                // Set the (button) interaction
                botCache.buttons.set(cmdFile.buttons[i].id, cmdFile.buttons[i].onClick);
            }
        }

        commands.push(cmdFile.command.data.toJSON());
    }

    // ========== SUBMIT COMMANDS ==========
    printLog('Started refreshing application (/) commands.', 'INFO');

    const url = config.devMode ?
        `https://discord.com/api/v8/applications/${config.botId}/guilds/${config.devGuild}/commands` :
        `https://discord.com/api/v8/applications/${config.botId}/commands`;

    const res = await fetch(url, {
        method: 'PUT',
        body: JSON.stringify(commands),
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bot ' + config.botToken
        }
    }).catch(ex => printLog(ex, 'ERROR'));

    printLog(`Application (/) commands PUT response: ${res.statusText} (${res.status})`, 'INFO');

    // ========== BOTLIST APIS ==========
    if (!config.devMode) {
        setInterval(async () => {
            const guildCount = (await manager.fetchClientValues('guilds.cache.size')).reduce((a, b) => a + b, 0);

            // =================================
            // TopGG (top.gg)
            await fetch(`https://top.gg/api/bots/${config.botId}/stats`, {
                method: 'POST',
                body: JSON.stringify({
                    server_count: guildCount,
                    shard_count: manager.totalShards,
                })
            });

            printLog(`Posted stats to TopGG (server_count: ${guildCount}, shard_count: ${manager.totalShards})`, 'INFO');

            // =================================
            // BotsForDiscord (discords.com/bots)
            await fetch(`https://discords.com/bots/api/bot/${config.botId}`, {
                method: 'POST',
                body: JSON.stringify({
                    server_count: guildCount
                }),
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': config.apis.discordsToken
                }
            });

            printLog(`Posted stats to BotsForDiscord (server_count: ${guildCount})`, 'INFO');

            // =================================
            // Discord Bots (discord.bots.gg)
            // TODO: This

            // =================================
            // Discord Bot List (discordbotlist.com)
            // TODO: This
        }, 1800000 /* 30 minutes */);
    }

    // ========== SHARD LAUNCH ==========
    manager.spawn().catch(console.error);
};

main();