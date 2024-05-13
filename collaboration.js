require('dotenv').config();
const WebSocket = require('ws');

const PORT = process.env.PORT || 8080;
let documentContent = "";

let documentUpdatesLog = [];

function broadcastDocument(clients, content) {
  clients.forEach(client => {
    if (client.readyState === WebSocket.OPEN) {
      sendDocumentToClient(client, content);
    }
  });
}

function sendDocumentToClient(client, content) {
  try {
    client.send(JSON.stringify({ type: 'document', content }));
  } catch (error) {
    console.error('Error sending document to a client:', error);
  }
}

function handleIncomingMessage(ws, message) {
  console.log('Received:', message);
  try {
    const data = JSON.parse(message);

    if (data.type === 'update') {
      documentContent = data.content;
      updateDocumentLog(data.content);
      broadcastDocument(wss.clients, documentContent);
    }
  } catch (error) {
    console.error('Error processing received message:', error);
  }
}

function updateDocumentLog(content) {
  const updateTimestamp = new Date().toISOString();
  documentUpdatesLog.push({ updateTimestamp, content });
  console.log(`Document updated at ${updateTimestamp}:`, content);
}

function getDocumentUpdateLogs() {
  return documentUpdatesLog;
}

const wss = new WebSocket.Server({ port: PORT });
console.log(`WebSocket server started on port: ${PORT}`);

wss.on('connection', (ws) => {
  console.log('Client connected');
  sendDocumentToClient(ws, documentContent); 

  ws.on('message', (message) => handleIncomingMessage(ws, message));

  ws.on('close', () => {
    console.log('Client disconnected');
  });
});

wss.on('error', (error) => {
  console.error('WebSocket server error:', error);
});