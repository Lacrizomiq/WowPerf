import Hero from "@/components/Home/Hero";
import TopPerformingSpec from "@/components/Home/TopPerformingSpec";
import Features from "@/components/Home/Features";
import CTA from "@/components/Home/CTA";
import Footer from "@/components/Home/Footer";

export default function Home() {
  return (
    <main className="bg-[#1A1D21]">
      {/* <Hero /> */}

      <div className="mt-8">
        <TopPerformingSpec />
      </div>
      <Features />
      <div className="p-4">
        <CTA />
      </div>
      <Footer />
    </main>
  );
}
