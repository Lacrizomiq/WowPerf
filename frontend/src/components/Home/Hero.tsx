import Image from "next/image";

export default function Hero() {
  return (
    <div className="relative h-[70vh] sm:h-[40vh] md:h-[50vh] lg:h-[60vh] xl:h-[70vh]">
      <Image
        src="/homepage.avif"
        alt="World of Warcraft Landscape"
        layout="fill"
        objectFit="cover"
        priority
      />
      <div className="absolute inset-0 bg-black bg-opacity-50 flex flex-col justify-center items-center text-center">
        <h1 className="text-5xl font-bold mb-4">World of Warcraft</h1>
        <p className="text-xl mb-8">
          Explore characters, equipment, and talents
        </p>
      </div>
    </div>
  );
}
