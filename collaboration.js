require('dotenv').config();
const WebSocket = require('ws');
const PORT = process.env.PORT || 8080;

let documentContent = "";

const broadcastDocument = (clients, content) => {
  clients.forEach(client => {
    try {
      if (client.readyState === WebSocket.OPEN) {
        client.send(JSON.stringify({ type: 'document', content }));
      }
    } catch (error) {
      console.error('Error broadcasting document to a client:', error);
    }
  });
};

const wss = new WebSocket.Server({ port: PORT });
console.log(`WebSocket server started on port: ${PORT}`);

wss.on('connection', (ws) => {
  console.log('Client connected');

  try {
    ws.send(JSON.stringify({ type: 'document', content: documentContent }));
  } catch (error) {
    console.error('Error sending document on connection:', error);
  }

  ws.on('message', (message) => {
    console.log('Received: %s', message);
    try {
      const data = JSON.parse(message);
  
      if (data.type === 'update') {
        documentContent = data.content;
        broadcastDocument(wss.clients, documentContent);
      }
    } catch (error) {
      console.error('Error processing received message:', error);
    }
  });

  ws.on('close', () => {
    console.log('Client disconnected');
  });
});

wss.on('error', (error) => {
  console.error('WebSocket server error:', error);
});