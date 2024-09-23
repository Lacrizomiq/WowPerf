import Image from "next/image";

export default function Hero() {
  return (
    <div className="relative h-[50vh] sm:h-[30vh] md:h-[40vh] lg:h-[40vh] xl:h-[50vh]">
      <div className="absolute inset-0 bg-[#002440] flex flex-col justify-center items-center text-center">
        <h1 className="text-5xl font-bold mb-4">Welcome to WowPerf</h1>
        <p className="text-xl mb-8">
          Explore characters, equipment, and talents to stay at the state of the
          art.
        </p>
      </div>
    </div>
  );
}
