// ================================
const config = require('../config')

// ================================
module.exports = async function (client) {
    // Set an interval for the activity so all guilds are loaded/cached before counting.
    setInterval(async function () {
        const guilds_result = await client.shard.fetchClientValues('guilds.cache.size')
        const guildCount = guilds_result.reduce((prev, count) => prev + count, 0)

        client.user.setActivity(`${guildCount} servers | ${client.shard.count} shards`, {
            type: 'WATCHING'
        })
    }, config.activityUpdateInterval * 60 * 1000)

    console.log('Fully started.')
}

module.exports.once = true