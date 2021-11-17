// ================================
const redis = require('redis')
const bluebird = require('bluebird')
const config = require('../config')

const { runQuery } = require('../structures/database')

// ================================
module.exports.botCache = {
    commands: new Map(),
    buttons: new Map()
}

// ================================
bluebird.promisifyAll(redis)
const redisClient = redis.createClient()
module.exports.redisClient = redisClient

module.exports.getFromRedis = async function (guildId) {
    if (await redisClient.existsAsync(guildId)) {
        return JSON.parse(await redisClient.getAsync(guildId))
    }

    return await cacheGuild(guildId)
}

module.exports.setInRedis = async function (guildId, data) {
    if (!await redisClient.existsAsync(guildId)) {
        await cacheGuild(guildId)
    }

    await redisClient.setAsync(guildId, JSON.stringify(data), 'EX', 60 * 60 * config.cacheExpireTime)
}

module.exports.removeFromRedis = async function (guildId) {
    if (await redisClient.existsAsync(guildId)) {
        await redisClient.delAsync(guildId)
    }
}

// ================================
async function cacheGuild(guildId) {
    // Get all the saved data from the guild
    let result = await runQuery('SELECT sug_channel, rep_channel, auto_approve, auto_reject, approve_emoji, reject_emoji, del_approved, del_rejected, blacklist FROM servers WHERE id = $1::text', [guildId])

    if (!result.rowCount) {
        // Register guild in database if it doesn't already exist
        await runQuery('INSERT INTO servers (id, premium) VALUES ($1::text, $2::bool)', [guildId, false])
        result.rows = [{ premium: false }]
    }

    await redisClient.setAsync(guildId, JSON.stringify(result.rows[0]), 'EX', 60 * 60 * config.cacheExpireTime)
    return result
}
