// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { MessageEmbed } = require('discord.js-light')
const config = require('../config')

// ================================
const data = new SlashCommandBuilder()
    .setName('list')
    .setDescription('Obtain a list of all open suggestions/reports.')
    .addStringOption(opt =>
        opt.setName('type')
            .setDescription('The specific list you want to see.')
            .addChoices([
                ['Suggestions', 'suggestions'],
                ['Reports', 'reports']
            ]).setRequired(true))

const execute = function (interaction) {
    const opt = interaction.options.getString('type')

    let desc = 'None'
    if (opt === 'suggestions') {
        // ...
    } else if (opt === 'reports') {
        // ...
    }

    const embed = new MessageEmbed()
        .setColor(config.embedColor.b)
        .setDescription(desc)

    interaction.reply({ embeds: [embed] })
}

// ================================
module.exports.command = {
    data: data,
    execute: execute
}