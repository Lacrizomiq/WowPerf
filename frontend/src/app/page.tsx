import Hero from "@/components/Home/Hero";
import Feature from "@/components/Home/Feature";
import Statistics from "@/components/Home/Statistics";
import Footer from "@/components/Home/Footer";

export default function Home() {
  return (
    <main className="bg-[#002440]">
      <Hero />
      <Feature />
      <Statistics />
      <Footer />
    </main>
  );
}
