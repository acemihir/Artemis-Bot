// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')

// ================================
const data = new SlashCommandBuilder()
	.setName('setstatus')
	.setDescription('Change the status of a suggestion/report/poll.')

	// Suggestions
	.addSubcommand(scmd =>
		scmd.setName('suggestions').setDescription('Change the status for a suggestion.')
			.addStringOption(opt => opt.setName('Suggestion ID').setDescription('The ID of the suggestion.'))
			.addStringOption(opt => opt.setName('Status').setDescription('The new status of the suggestion.')
				.addChoices([
					['Open', 'open'],
					['Approved', 'approved'],
					['Rejected', 'rejected'],
					['Considering', 'considering']
				]).setRequired(true))
	)

	// Reports
	.addSubcommand(scmd =>
		scmd.setName('reports').setDescription('Change the status for a report.')
			.addStringOption(opt => opt.setName('Report ID').setDescription('The ID of the report.'))
			.addStringOption(opt => opt.setName('Status').setDescription('The new status of the report.')
				.addChoices([
					['Open', 'open'],
					['Resolved', 'resolved'],
					['Progressing', 'progressing']
				]).setRequired(true))
	)

	// Polls
	.addSubcommand(scmd =>
		scmd.setName('polls').setDescription('Change the status for a poll.')
			.addStringOption(opt => opt.setName('Poll ID').setDescription('The ID of the poll.'))
			.addStringOption(opt => opt.setName('Status').setDescription('The new status of the poll.')
				.addChoices([
					['Open', 'open'],
					['Finished', 'finished']
				]).setRequired(true))
	)

const execute = async function(client, interaction) {
	await interaction.reply('to be done')
}

// ================================
module.exports.command = {
	isPremium: false,
	permLevel: 1,

	data: data,
	execute: execute
}