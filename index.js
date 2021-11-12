const { ShardingManager } = require('discord.js')
const config = require('./config')
const fetch = (...args) => import('node-fetch').then(({ default: fetch }) => fetch(...args))

const manager = new ShardingManager('app.js', { token: config.botToken })

// Spawn the shard
manager.on('shardCreate', shard => console.log(`Launched shard ${shard.id}`))
manager.spawn()


const main = async function () {

    // Production stuff
    if (!config.devMode) {
        // Every 30 minutes
        setInterval(async () => {
            // TopGG (top.gg)
            const guildCount = (await manager.fetchClientValues('guilds.cache.size')).reduce((a, b) => a + b, 0)

            await fetch(`https://top.gg/api/bots/${config.botId}/stats`, {
                method: 'POST',
                body: JSON.stringify({
                    server_count: guildCount,
                    shard_count: manager.totalShards,
                })
            })

            console.log(`Posted stats to TopGG (server_count: ${guildCount}, shard_count: ${manager.totalShards})`)

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
            })

            // Discord Bots (discord.bots.gg)
            // TODO: This

            // Discord Bot List (discordbotlist.com)
            // TODO: This
        }, 1800000)
    }
}

main()