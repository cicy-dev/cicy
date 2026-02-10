const express = require('express');
const app = express();
const port = process.env.PORT || 13001;
const fs = require('fs');
const path = require('path');

app.use(express.json({ limit: '50mb' }));

// MCP Protocol Version
const PROTOCOL_VERSION = "2024-11-05";

// Store messages and images
const messages = [];
const images = [];

// MCP Tools Definition
const tools = [
    {
        name: "send_message",
        description: "Send a text message to the server",
        inputSchema: {
            type: "object",
            properties: {
                message: {
                    type: "string",
                    description: "The message to send"
                }
            },
            required: ["message"]
        }
    },
    {
        name: "send_image",
        description: "Send an image to the server (base64 encoded)",
        inputSchema: {
            type: "object",
            properties: {
                name: {
                    type: "string",
                    description: "Image filename"
                },
                data: {
                    type: "string",
                    description: "Base64 encoded image data"
                },
                mimeType: {
                    type: "string",
                    description: "MIME type (e.g., image/png, image/jpeg)"
                }
            },
            required: ["name", "data"]
        }
    },
    {
        name: "get_messages",
        description: "Get all messages from the server",
        inputSchema: {
            type: "object",
            properties: {}
        }
    },
    {
        name: "get_images",
        description: "Get all images from the server",
        inputSchema: {
            type: "object",
            properties: {}
        }
    },
    {
        name: "clear_messages",
        description: "Clear all messages and images",
        inputSchema: {
            type: "object",
            properties: {}
        }
    }
];

// MCP Server Capabilities
const serverCapabilities = {
    tools: {
        listChanged: true
    },
    logging: {},
    resources: {}
};

// JSON-RPC Helper
function createResponse(id, result) {
    return {
        jsonrpc: "2.0",
        id,
        result
    };
}

function createError(id, code, message, data) {
    return {
        jsonrpc: "2.0",
        id,
        error: {
            code,
            message,
            data
        }
    };
}

// Convert image to ASCII art
function imageToAscii(base64Data, width = 80) {
    const chars = ' .:-=+*#%@';
    const lines = [];
    
    // Simple representation - in real implementation, decode and process image
    lines.push('┌' + '─'.repeat(width - 2) + '┐');
    lines.push('│' + ' '.repeat(width - 2) + '│');
    lines.push('│' + '  [Image Preview]'.padEnd(width - 2) + '│');
    lines.push('│' + ' '.repeat(width - 2) + '│');
    lines.push('└' + '─'.repeat(width - 2) + '┘');
    
    return lines.join('\n');
}

// MCP Methods
const mcpMethods = {
    initialize: (id, params) => {
        return createResponse(id, {
            protocolVersion: PROTOCOL_VERSION,
            capabilities: serverCapabilities,
            serverInfo: {
                name: "opencode-message-server",
                version: "1.0.0"
            }
        });
    },

    "tools/list": (id, params) => {
        return createResponse(id, {
            tools: tools
        });
    },

    "tools/call": (id, params) => {
        const { name, arguments: args } = params;
        
        switch (name) {
            case "send_message":
                if (!args || !args.message) {
                    return createError(id, -32602, "Invalid params: message required");
                }
                
                const message = {
                    type: 'text',
                    text: args.message,
                    timestamp: new Date().toISOString(),
                    id: messages.length + 1
                };
                messages.push(message);
                
                console.log(`[${new Date().toLocaleTimeString()}] Received message: ${args.message}`);
                
                // 随机长度回复
                const replies = [
                    '好',
                    '了解',
                    '收到了',
                    '明白了，我会处理',
                    '好的，我知道了，马上开始',
                    '收到你的消息了，我会认真处理这个请求',
                    '明白了，这个任务我会仔细完成，请放心',
                    '好的，我已经收到你的指示了，会按照要求来做',
                    '了解，我会立即开始处理这个任务，完成后会及时汇报结果'
                ];
                const randomReply = replies[Math.floor(Math.random() * replies.length)];
                
                return createResponse(id, {
                    content: [
                        {
                            type: "text",
                            text: randomReply
                        }
                    ],
                    isError: false
                });

            case "send_image":
                if (!args || !args.name || !args.data) {
                    return createError(id, -32602, "Invalid params: name and data required");
                }
                
                const image = {
                    type: 'image',
                    name: args.name,
                    mimeType: args.mimeType || 'image/png',
                    data: args.data.substring(0, 100) + '...', // Truncate for log
                    timestamp: new Date().toISOString(),
                    id: images.length + 1
                };
                images.push({
                    ...image,
                    data: args.data // Keep full data
                });
                
                console.log(`[${new Date().toLocaleTimeString()}] Received image: ${args.name}`);
                
                // Generate ASCII preview
                const asciiPreview = imageToAscii(args.data);
                
                return createResponse(id, {
                    content: [
                        {
                            type: "text",
                            text: `Image received: ${args.name}`
                        },
                        {
                            type: "image",
                            data: args.data,
                            mimeType: image.mimeType
                        }
                    ],
                    isError: false
                });

            case "get_messages":
                const allContent = [...messages, ...images.map(img => ({
                    type: 'image',
                    ...img
                }))];
                
                return createResponse(id, {
                    content: [
                        {
                            type: "text",
                            text: JSON.stringify(allContent, null, 2)
                        }
                    ],
                    isError: false
                });

            case "get_images":
                return createResponse(id, {
                    content: [
                        {
                            type: "text",
                            text: `Found ${images.length} images`
                        },
                        ...images.map(img => ({
                            type: "image",
                            data: img.data,
                            mimeType: img.mimeType
                        }))
                    ],
                    isError: false
                });

            case "clear_messages":
                messages.length = 0;
                images.length = 0;
                return createResponse(id, {
                    content: [
                        {
                            type: "text",
                            text: "All messages and images cleared"
                        }
                    ],
                    isError: false
                });

            default:
                return createError(id, -32601, `Method not found: ${name}`);
        }
    }
};

// Main JSON-RPC endpoint
app.post('/mcp', (req, res) => {
    const { jsonrpc, id, method, params } = req.body;
    
    if (jsonrpc !== "2.0") {
        return res.status(400).json(createError(id, -32600, "Invalid Request: jsonrpc must be 2.0"));
    }
    
    if (mcpMethods[method]) {
        const result = mcpMethods[method](id, params || {});
        res.json(result);
    } else {
        res.json(createError(id, -32601, `Method not found: ${method}`));
    }
});

// Legacy REST endpoint
// 随机回复句子
const responses = [
    '好',
    '收到',
    '了解',
    '明白',
    '知道了',
    '没问题',
    '好的好的',
    '收到了',
    '明白了',
    '我知道了',
    '好的我明白',
    '收到你的消息',
    '了解了解',
    '明白了会处理',
    '好的马上开始',
    '收到了我会认真处理',
    '明白了这个任务我会仔细完成',
    '好的我已经收到你的指示了',
    '了解我会立即开始处理这个任务',
    '收到了我会按照要求来做请放心',
    '明白了我会认真完成这个任务完成后会及时汇报结果'
];

function getRandomResponse() {
    return responses[Math.floor(Math.random() * responses.length)];
}

app.post('/message', (req, res) => {
    const { message } = req.body;
    if (message) {
        messages.push({
            type: 'text',
            text: message,
            timestamp: new Date().toISOString(),
            id: messages.length + 1
        });
        console.log(`[${new Date().toLocaleTimeString()}] Received: ${message}`);
        
        // 随机选择一个回复
        const reply = getRandomResponse();
        res.json({ success: true, message: reply });
    } else {
        res.status(400).json({ success: false, error: 'No message provided' });
    }
});

app.get('/messages', (req, res) => {
    res.json([...messages, ...images]);
});

// Health check
app.get('/health', (req, res) => {
    res.json({ 
        status: "ok", 
        protocol: "mcp", 
        version: PROTOCOL_VERSION,
        images: images.length,
        messages: messages.length
    });
});

app.listen(port, () => {
    console.log(`MCP Server running on http://localhost:${port}`);
    console.log(`JSON-RPC endpoint: POST http://localhost:${port}/mcp`);
    console.log(`Supports: text messages, images (base64)`);
});
