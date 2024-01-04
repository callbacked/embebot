import discord
import re
import os
import asyncio
import logging
from logging.handlers import RotatingFileHandler

async def suppress_embed(message):
        try:
            await message.edit(suppress=True)
            logging.info(f'Successfully suppressed embed') 
            return True
        except Exception as e:
            logging.warning(f'Failed to suppress embed')

async def wait_for_embed(message):
    logging.info(f'Starting to wait for embed for message ID: {message.id}')
    # await asyncio.sleep(0.5)
    updated_message = await message.channel.fetch_message(message.id)
    if updated_message.embeds:
        return True
    return False