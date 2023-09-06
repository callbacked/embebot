import discord
import re
import os
import asyncio
import logging
from logging.handlers import RotatingFileHandler

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

@client.event
async def on_ready():
    logging.info(f'{client.user} has connected to Discord!')

async def suppress_embed_with_retry(message, max_retries=3):
    for i in range(max_retries):
        await asyncio.sleep(0.5)  
        try:
            await message.edit(suppress=True)
            logging.info(f'Successfully suppressed embed preview on attempt {i + 1}') 
            return True
        except Exception as e:
            logging.warning(f'Failed to suppress embed preview on attempt {i + 1}: {e}')
    return False

async def wait_for_embed(message, timeout=5): ##changed from 10 to 1 because of discord vanilla embed failure for twitter on 6/30
    logging.info(f'Starting to wait for embed for message ID: {message.id}')
    for _ in range(timeout * 2):
        await asyncio.sleep(0.5)
        updated_message = await message.channel.fetch_message(message.id)
        if updated_message.embeds:
            return True
    return False
  
@client.event
async def on_message(message):
    if message.author == client.user:
        return

    twitter_match = re.search(r'https?://(?:www\.)?twitter\.com/[\w\-\_]+/status/\d+', message.content)
    x_match = re.search(r'https?://(?:www\.)?x\.com/[\w\-\_]+/status/\d+', message.content)
    tiktok_match = re.search(r'https?://(?:www\.)?tiktok\.com/.+', message.content)
    instagram_reel_match = re.search(r'https://www\.instagram\.com/reel/.*', message.content)
    
    # created a tuple to hold the regex link matches, base link, and embed in order to avoid ---
    # -- nested elif statements and create a single function
            # probably will just have this stored in a json file in the future

    match = [
    (twitter_match, 'twitter.com', 'vxtwitter.com'), 
    (x_match,'x.com', 'vxtwitter.com'),
    (instagram_reel_match, 'instagram.com', 'ddinstagram.com'),
    (tiktok_match, 'tiktok.com', 'vxtiktok.com')]

    for matched_link, base_link, vx_link in match:
        if matched_link:
            # suppresses the embed of the message containing the original matched link
            # to avoid clutter in chat 
            embed_found = await wait_for_embed(message)
            if embed_found:
                success = await suppress_embed_with_retry(message)
                if not success:
                    logging.warning("Failed to suppress embed preview after multiple attempts.")
            
            # generates the embed link based on the link it has found
            vx_url = matched_link.group(0).replace(base_link, vx_link)
            logging.info(f'Sending {vx_link} link: {vx_url}')
            await message.reply("[⠀]" + "(" + vx_url + ")", mention_author=False)
    
client.run(TOKEN)