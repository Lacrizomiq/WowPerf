import Image from "next/image";

export default function FeaturedContent() {
  const featuredItems = [
    {
      title: "Top Mythic + players",
      image: "/tww.png",
      description: "See the top players and mythic+ scores",
    },
    {
      title: "Top Guilds",
      image: "/worldfirst.png",
      description: "See the rankings of the best guilds",
    },
    {
      title: "Character Builds",
      image: "/build.png",
      description: "Discover popular character builds",
    },
  ];

  return (
    <section className="py-12 bg-gradient-dark">
      <div className="container mx-auto px-4">
        <h2 className="text-3xl font-bold mb-8 text-center text-gradient-glow">
          Featured Content
        </h2>
        <div className="flex flex-wrap justify-center gap-8">
          {featuredItems.map((item, index) => (
            <div
              key={index}
              className="w-full sm:w-[calc(50%-1rem)] lg:w-[calc(33.333%-1rem)] max-w-sm bg-deep-blue bg-opacity-50 rounded-lg overflow-hidden shadow-lg hover:scale-105 transition duration-300 glow-effect"
            >
              <Image
                src={item.image}
                alt={item.title}
                width={400}
                height={200}
                objectFit="cover"
              />
              <div className="p-6">
                <h3 className="text-xl font-semibold mb-2 text-gradient-glow">
                  {item.title}
                </h3>
                <p className="text-blue-200">{item.description}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
