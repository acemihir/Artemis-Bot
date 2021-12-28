const { ShardingManager } = require('discord.js-light');
const config = require('./config');
const { printLog } = require('./utils');
const fetch = (...args) => import('node-fetch').then(({ default: fetch }) => fetch(...args));

const manager = new ShardingManager('app.js', { token: config.botToken });

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

manager.on('shardCreate', shard => printLog('Launched shard.', 'INFO', shard.id));
manager.spawn().catch(console.error);