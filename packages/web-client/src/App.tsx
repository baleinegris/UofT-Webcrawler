import { useState } from 'react'
import Messages from './components/Messages';
import './App.css'

function App() {
  const [messages, setMessages] = useState<{ role: string; content: string; timestamp: string }[]>([]);

  async function handleSearch(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const rawQuery = form.query.value;
    
    // Clear the input field immediately after getting the value
    form.reset();
    
    const query = {
      query: rawQuery,
    }
    const CHATBOT_URL = import.meta.env.VITE_CHATBOT_URL || 'http://localhost:9001/query';
    const response = await fetch(CHATBOT_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(query),
    });
    if (!response.ok) {
      throw Error(`HTTP error ${response.status}`);
    }

    if (!response.body) {
      throw Error('Response body is null');
    }
    const reader = response.body.getReader();
    const decoder = new TextDecoder('utf-8');
    setMessages((prevMessages) => [
      ...prevMessages,
      { role: 'User', content: rawQuery, timestamp: new Date().toISOString() },
      { role: 'Agent', content: '', timestamp: new Date().toISOString() }, // Placeholder for streaming response
    ]);
    let fullResponse = '';
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      const chunk = decoder.decode(value, { stream: true });
      fullResponse += chunk;
      setMessages((prevMessages) => {
        const newMessages = [...prevMessages];
        newMessages[newMessages.length - 1] = {
          role: 'Agent',
          content: fullResponse,
          timestamp: newMessages[newMessages.length - 1].timestamp, // Keep the same timestamp
        };
        return newMessages;
      });
    }
  }
  return (
    <div className='w-[100vw] h-[100vh] absolute left-0 top-0 flex flex-col'>
      <div className="bg-blue-500 text-white p-4 rounded-lg w-full">
        <p className="text-lg font-bold">Welcome to SearchUofT!</p>
        <p>Search a question and get an answer scraped from the utoronto web.</p>
      </div>
      <div className="flex-grow overflow-y-auto flex justify-center items-start">
        <Messages messages={messages} />
      </div>
      <form onSubmit={handleSearch} className="mt-4 flex items-center fixed bottom-[5%] left-[50%] transform -translate-x-1/2 w-full max-w-4xl">
      <input
          type="text"
          placeholder="Search for a question..."
          className="flex-grow mr-2 p-2 rounded border w-full"
          name="query"
      />
      <button type='submit'className="bg-blue-500 text-white p-2 rounded">
        Search
      </button>
      </form>
    </div>
  )
}

export default App
