// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { createId } = require('../utils')

// ================================
const data = new SlashCommandBuilder()
	.setName('report')
	.setDescription('Create a report.')

const execute = async function(client, interaction) {
	await interaction.reply('to be done')
	const repId = createId('r_')
}

// ================================
module.exports.command = {
	isPremium: false,
	permLevel: 0,

	data: data,
	execute: execute
}