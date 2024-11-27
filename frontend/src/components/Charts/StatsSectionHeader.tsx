// components/Charts/StatsSectionHeader.tsx
interface StatsSectionHeaderProps {
  title: string;
  total: number;
  subtitle?: string;
}

const StatsSectionHeader: React.FC<StatsSectionHeaderProps> = ({
  title,
  total,
  subtitle,
}) => (
  <div className="mb-4">
    <h3 className="text-xl font-bold text-white">
      {title} - Total: {total.toLocaleString()} players
    </h3>
    {subtitle && <p className="text-gray-300 mt-1">{subtitle}</p>}
  </div>
);

export default StatsSectionHeader;
