// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')

// ================================
const data = new SlashCommandBuilder()
	.setName('help')
	.setDescription('Obtain information general information about the bot.')
	.addStringOption(opt =>
	opt.setName('command')
		.setDescription('Obtain information about a certain command.')
		.addChoices([
			['Suggest', 'suggest'],
			['Report', 'report'],
			['About', 'about'],
			['Help', 'help'],
			['Graph', 'graph'],
			['SetStatus', 'setstatus'],
			['Move', 'move'],
			['Poll', 'poll'],
			['Setup', 'setup'],
			['Config', 'config'],
			['Blacklist', 'blacklist']
		]).setRequired(true))

const execute = async function(client, interaction) {
		await interaction.reply(
`\`\`\`asciidoc
== Commands ==
[View the autocompletion for more detailed explanation.]
= User =
/suggest    :: Create a suggestion for this server.
/report     :: Create a report for this server.
/about      :: Obtain information about this bot.
/help       :: Shows this particular help message.
/graph      :: Get a graph displaying suggestion/report information. (Premium Only)
= Staff =
/setstatus  :: Change the status of a Suggestion or Report.
/move       :: Move a suggestion or report to another channel.
/poll       :: Create a poll in this server.
= Admin =
/setup      :: Setup the bot. (This will set the staffrole and suggestion channel)
/config     :: Modify the bot's configuration for your server.
/blacklist  :: Prevent someone from creating suggestions and/or reports.
\`\`\``
		)
}

// ================================
module.exports.command = {
	isPremium: false,
	permLevel: 0,

	data: data,
	execute: execute
}