const { logger } = require('../utils')

module.exports = {
    once: true,
    execute(client) {
        // Set an interval for the activity so all guilds are loaded/cached before counting.
        setInterval(async function () {
            const guilds_result = await client.shard.fetchClientValues('guilds.cache.size')
            const guildCount = guilds_result.reduce((prev, count) => prev + count, 0)

            await client.user.setActivity(`${guildCount} servers | ${client.shard.count} shards`, {
                type: 'WATCHING'
            })
        }, 15 * 60 * 1000)

        logger.debug(client.shard.ids + ' Fully started.')
    },
}