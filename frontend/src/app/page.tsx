import Hero from "@/components/Home/Hero";
import Feature from "@/components/Home/Feature";
import Statistics from "@/components/Home/Statistics";
import CTA from "@/components/Home/CTA";
import Footer from "@/components/Home/Footer";

export default function Home() {
  return (
    <main className="bg-[#1A1D21]">
      <Hero />
      {/* Add a negative margin here */}
      <div className="-mt-8">
        <Feature />
      </div>
      <Statistics />
      <CTA />
      <Footer />
    </main>
  );
}
