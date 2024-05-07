require('dotenv').config();
const WebSocket = require('ws');
const PORT = process.env.PORT || 8080;

let documentContent = "";

const broadcastDocument = (clients, content) => {
  clients.forEach(client => {
    if (client.readyState === WebSocket.OPEN) {
      client.send(JSON.stringify({ type: 'document', content }));
    }
  });
};

const wss = new WebSocket.Server({ port: PORT });
console.log(`WebSocket server started on port: ${PORT}`);

wss.on('connection', (ws) => {
  console.log('Client connected');
  ws.send(JSON.stringify({ type: 'document', content: documentContent }));

  ws.on('message', (message) => {
    console.log('Received: %s', message);
    const data = JSON.parse(message);
    if (data.type === 'update') {
      documentContent = data.content;
      broadcastDocument(wss.clients, documentContent);
    }
  });

  ws.on('close', () => {
    console.log('Client disconnected');
  });
});

wss.on('error', (error) => {
  console.error('WebSocket server error:', error);
});