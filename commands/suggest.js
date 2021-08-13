// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')

// ================================
const data = new SlashCommandBuilder()
	.setName('suggest')
	.setDescription('Create a suggestion.')

const execute = async function(client, interaction) {
	await interaction.reply('to be done')
	const sugId = createId('s_')
}

// ================================
module.exports.command = {
	isPremium: false,
	permLevel: 0,

	data: data,
	execute: execute
}