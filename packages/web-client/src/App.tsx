import './App.css'

function App() {

  return (
    <>
      <div className="bg-blue-500 text-white p-4 rounded-lg">
        <p className="text-lg font-bold">Welcome to SearchUofT!</p>
        <p>Search a question and get an answer scraped from the utoronto web.</p>
      </div>
      <input
          type="text"
          placeholder="Search for a question..."
          className="flex-grow mr-2 p-2 rounded border border-gray-300 w-full"
      />
    </>
  )
}

export default App
