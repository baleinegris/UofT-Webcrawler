import './App.css'

function App() {

  async function handleSearch(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const rawQuery = event.currentTarget.query.value;
    const query = {
      query: rawQuery,
    }
    const response = await fetch('http://localhost:8000/query', {
      method: 'GET',
      body: JSON.stringify(query),
    });
    if (!response.ok) {
      console.error('Error fetching data:', response.statusText);
      return;
    }
    const data = await response.json();
    console.log('Response from server:', data);
  }
  return (
    <>
      <div className="bg-blue-500 text-white p-4 rounded-lg">
        <p className="text-lg font-bold">Welcome to SearchUofT!</p>
        <p>Search a question and get an answer scraped from the utoronto web.</p>
      </div>
      <form onSubmit={handleSearch} className="mt-4 flex items-center">
      <input
          type="text"
          placeholder="Search for a question..."
          className="flex-grow mr-2 p-2 rounded border border-gray-300 w-full"
          name="query"
      />
      <button type='submit'className="bg-blue-500 text-white p-2 rounded">
        Search
      </button>
      </form>
    </>
  )
}

export default App
