// ================================
const redis = require('redis')
const bluebird = require('bluebird')
const config = require('../config')

const { runQuery } = require('../structures/database')

// ================================
module.exports.botCache = {
	commands: new Map(),
	buttonActions: new Map()
}

// ================================
bluebird.promisifyAll(redis)
const redisClient = redis.createClient()
module.exports.redisClient = redisClient

module.exports.getFromRedis = async function(guildId) {
	if (await redisClient.existsAsync(guildId)) {
		return JSON.parse(await redisClient.getAsync(guildId))
	}

	return await cacheGuild(guildId)
}

module.exports.setInRedis = async function(guildId, data) {
	if (!await redisClient.existsAsync(guildId)) {
		await cacheGuild(guildId)
	}

	await redisClient.setAsync(guildId, JSON.stringify(newData), 'EX', 60 * 60 * config.cacheExpireTime)
}

module.exports.removeFromRedis = async function(guildId) {
	if (await redisClient.existsAsync(guildId)) {
		await redisClient.delAsync(guildId)
	}
}

// ================================
async function cacheGuild(guildId) {
	let result = await runQuery('SELECT staff_role, approve_emoji, reject_emoji, premium FROM servers WHERE id = $1::text', [guildId])

	if (!result.rowCount) {
		// Register guild in database if it doesn't already exist
		await runQuery('INSERT INTO servers (id, premium) VALUES ($1::text, $2::bool)', [guildId, false])
		result.rows = [{ premium: false }]
	}

	const data = {
		staffRole: result.rows[0].staff_role,
		autoApprove: result.rows[0].auto_approve || -1,
		autoReject: result.rows[0].auto_reject || -1,
		isPremium: result.rows[0].premium
	}

	await redisClient.setAsync(guildId, JSON.stringify(data), 'EX', 60 * 60 * config.cacheExpireTime)
	return data
}
