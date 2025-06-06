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
      <div className="-mt-2">
        <Feature />
      </div>
      <Statistics />
      <div className="p-4">
        <CTA />
      </div>
      <Footer />
    </main>
  );
}
