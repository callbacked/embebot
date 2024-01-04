import discord
import re
import os
import logging
import json
import configparser
from logging.handlers import RotatingFileHandler
from suppression import suppress_embed 
from suppression import wait_for_embed

TOKEN = os.getenv('DISCORD_BOT_TOKEN')

# logging
log_formatter = logging.Formatter('%(asctime)s %(levelname)s %(message)s')
console_handler = logging.StreamHandler()
console_handler.setFormatter(log_formatter)

file_handler = RotatingFileHandler('bot.log', maxBytes=50*1024*1024, backupCount=1)
file_handler.setFormatter(log_formatter)

logging.basicConfig(
    level=logging.INFO,
    handlers=[
        console_handler,
        file_handler
    ]
)

intents = discord.Intents.default()
intents.message_content = True
client = discord.Client(intents=intents)

with open('match.json', 'r') as f:
    match_config = json.load(f)

config = configparser.ConfigParser()
config.read("config.ini")

if config['Settings'].getboolean('EndpointOverride'):
    logging.info(f'Endpoint overrides found, ignoring defaults')
    for item in match_config:
        service_name = item["base_link"].replace(".com", "")  
        vx_link_override_key = f"{service_name}_vx_link"
        if vx_link_override_key in config['vx_links'] and config['vx_links'][vx_link_override_key].lower() != 'default':
            item["vx_link"] = config['vx_links'][vx_link_override_key]


@client.event
async def on_ready():
    logging.info(f'{client.user} has connected to Discord!')
    await client.change_presence(status=discord.Status.dnd, activity=discord.Activity(type=discord.ActivityType.watching, name="your messages"))

    
@client.event
async def on_message(message):
    if message.author == client.user:
        return
    
    for item in match_config:
        pattern = re.compile(item["pattern"])
        matched_link = pattern.search(message.content)

        if matched_link:
            embed_found = await wait_for_embed(message)
            if embed_found:
                success = await suppress_embed(message)
                if not success:
                    logging.warning("Failed to suppress embed preview after multiple attempts.")

            vx_url = matched_link.group(0).replace(item["base_link"], item["vx_link"])
            logging.info(f'Sending {item["vx_link"]} link: {vx_url}')
            await message.reply("[⠀]" + "(" + vx_url + ")", mention_author=False)

    
client.run(TOKEN)