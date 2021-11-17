// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')

// ================================
const data = new SlashCommandBuilder()
    .setName('help')
    .setDescription('Obtain information general information about the bot.')
    .addStringOption(opt => opt.setName('command')
        .setDescription('Obtain information about a certain command.')
        .addChoices([
            ['Suggest', 'suggest'],
            ['Report', 'report'],
            ['About', 'about'],
            ['Help', 'help'],
            ['List', 'list'],
            ['Graph', 'graph'],
            ['SetStatus', 'setstatus'],
            ['Move', 'move'],
            ['Blacklist', 'blacklist'],
            ['Setup', 'setup'],
            ['Config', 'config']
        ]).setRequired(false))

const execute = function (interaction) {
    let desc = `\`\`\`asciidoc
== Commands ==
[View the autocompletion for more detailed explanation.]
= User =
/suggest    :: Create a suggestion for this server.
/report     :: Create a report for this server.
/about      :: Obtain information about this bot.
/help       :: Shows this particular help message.
/list       :: Get a list of suggestions, staffmembers can list reports as well (Premium Only)
/graph      :: Get a graph displaying suggestion/report information. (Premium Only)
= Staff =
/setstatus  :: Change the status of a Suggestion or Report.
/move       :: Move a suggestion or report to another channel.
/blacklist  :: Prevent someone from creating suggestions and/or reports. (Premium Only)
= Admin =
/setup      :: Setup the bot. (This will set the staffrole and suggestion channel)
/config     :: Modify the bot's configuration for your server.\`\`\``

    const opt = interaction.options.getString('command')

    if (opt != null) {
        if (opt === 'suggest') {
            desc = 'The suggest command can be used by people to suggest their ideas.'
        } else if (opt === 'report') {
            desc = 'The report command can be used by users to send a report to the staffmembers.'
        } else if (opt === 'about') {
            desc = 'This particular command shows some information about this Discord bot.'
        } else if (opt === 'help') {
            desc = 'The command that can be used to obtain a list with all commands + explanations.'
        } else if (opt === 'list') {
            desc = 'A premium command that can be used to display all active suggestions/reports.'
        } else if (opt === 'graph') {
            desc = 'A premium command that can be used to display some information in graphs.'
        } else if (opt === 'setstatus') {
            desc = 'A command that can be used by staffmembers to change the status of a suggestion/report.'
        } else if (opt === 'move') {
            desc = 'This command can be used to move the message of a suggestion/report to another channel.'
        } else if (opt === 'setup') {
            desc = 'This admin only command can be used to conifgure all fundamentally required settings for this bot.'
        } else if (opt === 'config') {
            desc = 'Configuring the bot can be done using this command. This command is admin only.'
        } else if (opt === 'blacklist') {
            desc = 'Users can be blacklisted from submitting suggestions/reports after they\'ve been added to the blacklist. This command is premium only.'
        }
    }

    interaction.reply(desc)
}

// ================================
module.exports.command = {
    data: data,
    execute: execute
}