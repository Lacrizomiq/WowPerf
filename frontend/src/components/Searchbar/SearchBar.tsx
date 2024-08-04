export default function SearchBar() {
  return (
    <div className="bg-gray-800 py-8">
      <div className="container mx-auto">
        <form className="flex justify-center">
          <input
            type="text"
            placeholder="Search for a character..."
            className="w-1/2 px-4 py-2 rounded-l-lg focus:outline-none focus:ring-2 focus:ring-blue-400 text-black"
          />
          <button
            type="submit"
            className="bg-blue-500 text-white px-6 py-2 rounded-r-lg hover:bg-blue-600 transition duration-300"
          >
            Search
          </button>
        </form>
      </div>
    </div>
  );
}
