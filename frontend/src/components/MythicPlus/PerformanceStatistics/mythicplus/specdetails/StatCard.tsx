import React from "react";
import { ArrowUp, ArrowDown } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
interface StatCardProps {
  title: string;
  value: string;

  trend?: "up" | "down";
  isComing?: boolean;
}

const StatCard: React.FC<StatCardProps> = ({
  title,
  value,
  trend,
  isComing = false,
}) => {
  return (
    <Card className="bg-slate-800/30 border-slate-700">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-slate-300">
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center">
          {isComing ? (
            <Badge className="ml-2 bg-purple-600 text-[10px]">Soon</Badge>
          ) : (
            <>
              <div className="text-xl font-bold">{value}</div>
              {trend === "up" && (
                <ArrowUp className="ml-2 h-3 w-3 text-green-400" />
              )}
              {trend === "down" && (
                <ArrowDown className="ml-2 h-3 w-3 text-red-400" />
              )}
            </>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

export default StatCard;
