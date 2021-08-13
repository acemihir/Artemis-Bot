const { ShardingManager } = require('discord.js')
const config = require('./config')
const { AutoPoster } = require('topgg-autoposter')

const manager = new ShardingManager('app.js', { token: config.botToken })

if (!config.devMode) {
    const ap = AutoPoster(config.apis.topggToken, manager)
    ap.on('posted', (stats) => {
        console.log(`Posted stats to Top.gg | ${stats.serverCount} servers`)
    })
}

manager.on('shardCreate', shard => console.log(`Launched shard ${shard.id}`))
manager.spawn()