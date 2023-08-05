import { NextApiRequest, NextApiResponse } from "next";
import Redis from 'ioredis';

const handler = async (req, res) => {
  res.setHeader('Cache-Control', 'no-cache');
  res.setHeader('Content-Type', 'text/event-stream');
  res.setHeader('Connection', 'keep-alive');
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Credentials', 'true');
  res.setHeader('Access-Control-Allow-Methods', 'GET, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  const redisConfig = {
    host: 'localhost',
    port: 6379,
  };
  const redis = new Redis(redisConfig);

  // Subscribe to a channel
  await redis.subscribe('thumbnail');

  // Listen for incoming messages
  redis.on('message', (channel, message) => {
    console.log(`Received message in channel ${channel}: ${message}`);
    res.write(`data: ${message}\n\n`);
  });

  // Keep the connection open
  req.on('close', () => {
    redis.unsubscribe('thumbnail');
    redis.quit();
  });
};

export default handler;
