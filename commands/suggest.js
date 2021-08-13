// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { createId, filterText } = require('../utils')

// ================================
const data = new SlashCommandBuilder()
	.setName('suggest')
	.setDescription('Create a suggestion.')

const execute = async function(client, interaction) {
	const input = await interaction.fetchReply()
	print(input)

	await interaction.reply('to be done')
	const sugId = createId('s_')
}


module.exports.buttonActions = [
	{
		id = 'sug_upvote',
		onClick = async function(client, interaction) {
			// Update the message and update the count
		}
	},
	{
		id = 'sug_downvote',
		onClick = async function(client, interaction) {
			// Update the message and update the count
		}
	}
]
 

// ================================
module.exports.command = {
	isPremium: false,
	permLevel: 0,

	data: data,
	execute: execute,
	buttons: buttonActions
}