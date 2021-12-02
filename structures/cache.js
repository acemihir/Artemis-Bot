// ================================
const bluebird = require('bluebird')
const redis = require('redis')
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
    let result = await runQuery('SELECT staff_role, sug_channel, rep_channel, auto_consider, auto_approve, auto_reject, approve_emoji, reject_emoji, del_approved, del_rejected, blacklist FROM servers WHERE id = $1::text', [guildId])

    if (!result.rowCount) {
        // Register guild in database if it doesn't already exist
        await runQuery('INSERT INTO servers (id, premium) VALUES ($1::text, $2::bool)', [guildId, false])
        result.rows = [{ premium: false }]
    }

    // Configure some default params
    const data = {
        staff_role: result.rows[0].staff_role,
        
        sug_channel: result.rows[0].sug_channel,
        rep_channel: result.rows[0].rep_channel,
        
        auto_consider: result.rows[0].auto_consider || -1,
        auto_approve: result.rows[0].auto_approve || -1,
        auto_reject: result.rows[0].auto_reject || -1,
        
        approve_emoji: result.rows[0].approve_emoji || '⬆️',
        reject_emoji: result.rows[0].reject_emoji || '⬇️',
        
        del_approved: result.rows[0].del_approved || false,
        del_rejected: result.rows[0].del_approved || false,
        
        blacklist: result.rows[0].del_approved || '[]',
    }

    // Set the data in the cache
    await redisClient.setAsync(guildId, JSON.stringify(data), 'EX', 60 * 60 * config.cacheExpireTime)
    return result
}
