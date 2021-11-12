// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { MessageEmbed, MessageActionRow, MessageButton } = require('discord.js')
const config = require('../config')
const { setInRedis, getFromRedis, botCache } = require('../structures/cache')
const { runQuery } = require('../structures/database')

// ================================
const data = new SlashCommandBuilder()
    .setName('config')
    .setDescription('Configure the bot to have it fit your needs.')

const execute = function (client, interaction) {
    if (!interaction.member.permissions.has('ADMINISTRATOR')) {
        return interaction.reply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.r)
                .setDescription('🟥 You need to have the `ADMINISTRATOR` permission to do that.')
            ]
        })
    }

    const embed = new MessageEmbed()
        .setAuthor('Config -> Main', client.user.avatarURL())
        .setColor(config.embedColor.b)
        .setDescription('Welcome to the config menu of Suggestions, you can use the buttons below to navigate your way through ' +
            'the options that can be configured. If you encounter a problem or a question raises then don\'t hesitate to ask our ' +
            'support team over in our [Support Server](https://discord.gg/3SYg3M5).')

    const row = new MessageActionRow().addComponents(
        new MessageButton()
            .setCustomId('conf_channels')
            .setLabel('Channels')
            .setStyle('PRIMARY')
            .setEmoji('#️⃣'),
        new MessageButton()
            .setCustomId('conf_roles')
            .setLabel('Roles')
            .setStyle('SUCCESS')
            .setEmoji('👥'),
        new MessageButton()
            .setCustomId('conf_behaviour')
            .setLabel('Behaviour')
            .setStyle('DANGER')
            .setEmoji('🦿')
        // Add an extra ("other") category here when needed
    )

    interaction.reply({ embeds: [embed], components: [row] })
}

module.exports.buttons = [
    // ================================
    // Channel Configuration
    {
        id: 'conf_channels',
        onClick: async function (client, interaction) {
            if (!interaction.member.permissions.has('ADMINISTRATOR')) {
                return interaction.reply({
                    embeds: [new MessageEmbed()
                        .setColor(config.embedColor.r)
                        .setDescription('🟥 You need to have the `ADMINISTRATOR` permission to do that.')
                    ], ephemeral: true
                })
            }

            // Create a new embed
            const embed = new MessageEmbed()
                .setAuthor('Config -> Channels', client.user.avatarURL())
                .setColor(config.embedColor.b)
                .setDescription('Click on the button with the label corresponding to the function for which you want to change the channel.')

            const row = new MessageActionRow().addComponents(
                new MessageButton()
                    .setCustomId('conf_channels_sug')
                    .setLabel('Suggestions')
                    .setStyle('SUCCESS')
                    .setEmoji('❔'),
                new MessageButton()
                    .setCustomId('conf_channels_rep')
                    .setLabel('Reports')
                    .setStyle('DANGER')
                    .setEmoji('❗'),
            )

            // Edit the message
            await interaction.message.edit({ embeds: [embed], components: [row] })
            interaction.deferUpdate()
        }
    },
    {
        id: 'conf_channels_sug',
        onClick: async function (client, interaction) {
            if (!interaction.member.permissions.has('ADMINISTRATOR')) {
                return interaction.reply({
                    embeds: [new MessageEmbed()
                        .setColor(config.embedColor.r)
                        .setDescription('🟥 You need to have the `ADMINISTRATOR` permission to do that.')
                    ], ephemeral: true
                })
            }

            // Create a new embed
            const embed = new MessageEmbed()
                .setAuthor('Config -> Channels -> Suggestions', client.user.avatarURL())
                .setColor(config.embedColor.b)
                .setDescription('In which channel should the suggestions show up? (Type: #channel)')

            interaction.reply({ embeds: [embed] })

            const filter = msg => msg.author.id === interaction.user.id

            // Get the channel
            const chnAwait = await interaction.channel.awaitMessages({ filter, max: 1, time: 25000, errors: ['time'] })
            await chnAwait.first().delete()

            const chnId = chnAwait.first().content.replace('<#', '').replace('>', '')
            if (interaction.guild.channels.cache.get(chnId) == null) {
                embed.setColor(config.embedColor.r)
                embed.setDescription('🟥 That\'s not a valid channel, please run the command again.')
                return interaction.editReply({ embeds: [embed] })
            }

            const currentCache = await getFromRedis(interaction.guildId)
            currentCache['sug_channel'] = chnId
            await setInRedis(interaction.guildId, currentCache)

            await runQuery('UPDATE servers SET sug_channel = $1::text WHERE id = $2::text', [chnId, interaction.guildId])

            embed.setColor(config.embedColor.g)
            embed.setDescription(`Suggestions will now show up in <#${chnId}>.`)

            await interaction.editReply({ embeds: [embed] })
            interaction.deferUpdate()
        }
    },
    {
        id: 'conf_channels_rep',
        onClick: async function (client, interaction) {
            if (!interaction.member.permissions.has('ADMINISTRATOR')) {
                return interaction.reply({
                    embeds: [new MessageEmbed()
                        .setColor(config.embedColor.r)
                        .setDescription('🟥 You need to have the `ADMINISTRATOR` permission to do that.')
                    ], ephemeral: true
                })
            }

            // Create a new embed
            const embed = new MessageEmbed()
                .setAuthor('Config -> Channels -> Reports', client.user.avatarURL())
                .setColor(config.embedColor.b)
                .setDescription('In which channel should the reports show up? (Type: #channel)')

            interaction.reply({ embeds: [embed] })

            const filter = msg => msg.author.id === interaction.user.id

            // Get the channel
            const chnAwait = await interaction.channel.awaitMessages({ filter, max: 1, time: 25000, errors: ['time'] })
            await chnAwait.first().delete()

            const chnId = chnAwait.first().content.replace('<#', '').replace('>', '')
            if (interaction.guild.channels.cache.get(chnId) == null) {
                embed.setColor(config.embedColor.r)
                embed.setDescription('🟥 That\'s not a valid channel, please run the command again.')
                return interaction.editReply({ embeds: [embed] })
            }

            const currentCache = await getFromRedis(interaction.guildId)
            currentCache['rep_channel'] = chnId
            await setInRedis(interaction.guildId, currentCache)

            await runQuery('UPDATE servers SET rep_channel = $1::text WHERE id = $2::text', [chnId, interaction.guildId])

            embed.setColor(config.embedColor.g)
            embed.setDescription(`Reports will now show up in <#${chnId}>.`)

            await interaction.editReply({ embeds: [embed] })
        }
    },

    // ================================
    // Role configuration
    {
        id: 'conf_roles',
        onClick: async function (client, interaction) {
            if (!interaction.member.permissions.has('ADMINISTRATOR')) {
                return interaction.reply({
                    embeds: [new MessageEmbed()
                        .setColor(config.embedColor.r)
                        .setDescription('🟥 You need to have the `ADMINISTRATOR` permission to do that.')
                    ], ephemeral: true
                })
            }

            // Create a new embed
            const embed = new MessageEmbed()
                .setAuthor('Config -> Roles', client.user.avatarURL())
                .setColor(config.embedColor.b)
                .setDescription('Click on the button with the label corresponding to the role you want to change.')

            const row = new MessageActionRow().addComponents(
                new MessageButton()
                    .setCustomId('conf_roles_staff')
                    .setLabel('Staff')
                    .setStyle('SUCCESS')
                    .setEmoji('🚧'),
            )

            // Edit the message
            await interaction.message.edit({ embeds: [embed], components: [row] })
            await interaction.deferUpdate()
        }
    },
    {
        id: 'conf_roles_staff',
        onClick: async function (client, interaction) {
            if (!interaction.member.permissions.has('ADMINISTRATOR')) {
                return interaction.reply({
                    embeds: [new MessageEmbed()
                        .setColor(config.embedColor.r)
                        .setDescription('🟥 You need to have the `ADMINISTRATOR` permission to do that.')
                    ], ephemeral: true
                })
            }

            // Create a new embed
            const embed = new MessageEmbed()
                .setAuthor('Config -> Roles -> Staff', client.user.avatarURL())
                .setColor(config.embedColor.b)
                .setDescription('Which role should be able to interact with created Suggestions & Reports? (Type: @role)')

            interaction.reply({ embeds: [embed] })

            const filter = msg => msg.author.id === interaction.user.id

            const roleAwait = await interaction.channel.awaitMessages({ filter, max: 1, time: 25000, errors: ['time'] })
            await roleAwait.first().delete()

            // The role directly to check and so we can get the role name later on
            const role = interaction.guild.roles.cache.get(roleAwait.first().content.replace('<@&', '').replace('>', ''))
            if (role == null) {
                embed.setColor(config.embedColor.r)
                embed.setDescription('🟥 That\'s not a valid role, please run the command again.')
                return interaction.editReply({ embeds: [embed] })
            }

            // Update the discord command permission
            const commands = await interaction.guild.commands.fetch()
            const permissions = []
            // Loop through all the commands
            for (const [key, value] of commands.entries()) {
                if (value.applicationId == interaction.applicationId) {
                    if (botCache.privCommands.includes(value.name)) {
                        permissions.push({
                            id: key,
                            permissions: [{
                                id: role.id,
                                type: 'ROLE',
                                permission: true,
                            }]
                        })
                    }
                }
            }
            // Set the actual permission
            await interaction.guild.commands.permissions.set({ fullPermissions: permissions })

            embed.setColor(config.embedColor.g)
            embed.setDescription(`The ${role.name} can now interact with Suggestions & Reports.`)

            await interaction.editReply({ embeds: [embed] })
        }
    },

    // ================================
    // Behaviour configuration
    {
        id: 'conf_behaviour',
        onClick: async function (client, interaction) {
            if (!interaction.member.permissions.has('ADMINISTRATOR')) {
                return interaction.reply({
                    embeds: [new MessageEmbed()
                        .setColor(config.embedColor.r)
                        .setDescription('🟥 You need to have the `ADMINISTRATOR` permission to do that.')
                    ], ephemeral: true
                })
            }

            // Create a new embed
            const embed = new MessageEmbed()
                .setAuthor('Config -> Channels', client.user.avatarURL())
                .setColor(config.embedColor.b)
                .setDescription('This section will be available soon.')

            // Edit the message
            await interaction.message.edit({ embeds: [embed] })
            interaction.deferUpdate()
        }
    }
]

// ================================
module.exports.command = {
    isPremium: false,
    privileged: true,

    data: data,
    execute: execute
}