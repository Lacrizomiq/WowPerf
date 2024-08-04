import Header from "@/components/Header/Header";
import Hero from "@/components/Home/Hero";
import SearchBar from "@/components/Searchbar/SearchBar";
import FeaturedContent from "@/components/Home/FeaturedContent";

export default function Home() {
  return (
    <main>
      <Header />
      <Hero />
      <SearchBar />
      <FeaturedContent />
    </main>
  );
}
