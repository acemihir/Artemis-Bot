// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { MessageActionRow, MessageEmbed, MessageButton } = require('discord.js-light')
const config = require('../config')

// ================================
const data = new SlashCommandBuilder()
    .setName('about')
    .setDescription('Obtain general information about the bot.')

const execute = function(_client, interaction) {
    const row = new MessageActionRow().addComponents(
        new MessageButton()
            .setURL('https://top.gg/bot/566616056165302282/invite/')
            .setLabel('Invite')
            .setEmoji('🤖')
            .setStyle('LINK'),
        new MessageButton()
            .setURL('https://github.com/jerskisnow/Suggestions/wiki/Donating')
            .setLabel('Donate')
            .setEmoji('💰')
            .setStyle('LINK'),
        new MessageButton()
            .setURL('https://top.gg/bot/566616056165302282/vote')
            .setLabel('Vote')
            .setEmoji('📰')
            .setStyle('LINK'),
        new MessageButton()
            .setURL('https://discord.gg/3SYg3M5')
            .setLabel('Discord')
            .setEmoji('👥')
            .setStyle('LINK')
    )

    const embed = new MessageEmbed()
    embed.setColor(config.embedColor.b)
    embed.setTitle('About Suggestions')
    embed.setDescription('Suggestions is a discord bot, created by CodedSnow (jerskisnow), that allows for perfect collaboration between members and staff members.' +
        ' Members can submit their ideas, a staff member can then approve, consider or reject them.' +
        ' Suggestions also offers support for reports, polls and a lot more!\n\nConsider taking a look at the buttons below this message for more information and how to support us.')

    interaction.reply({ embeds: [embed], components: [row] })
}

// ================================
module.exports.command = {
    isPremium: false,

    data: data,
    execute: execute
}