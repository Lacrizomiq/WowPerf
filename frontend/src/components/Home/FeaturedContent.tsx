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
    <section className="py-12 bg-gray-900">
      <div className="container mx-auto">
        <h2 className="text-3xl font-bold mb-8 text-center">
          Featured Content
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {featuredItems.map((item, index) => (
            <div
              key={index}
              className="bg-gray-800 rounded-lg overflow-hidden shadow-lg hover:scale-105 transition duration-300"
            >
              <Image
                src={item.image}
                alt={item.title}
                width={400}
                height={200}
                objectFit="cover"
              />
              <div className="p-6">
                <h3 className="text-xl font-semibold mb-2">{item.title}</h3>
                <p className="text-gray-400">{item.description}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
