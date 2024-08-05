import Header from "@/components/Header/Header";
import Hero from "@/components/Home/Hero";
import SearchBar from "@/components/Home/SearchBar";
import FeaturedContent from "@/components/Home/FeaturedContent";

export default function Home() {
  return (
    <main className="bg-gradient-dark">
      <Header />
      <Hero />
      <SearchBar />
      <FeaturedContent />
    </main>
  );
}
