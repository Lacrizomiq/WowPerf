import React from "react";
import { Info } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { Badge } from "@/components/ui/badge";
const ScoreEvolution: React.FC = () => {
  return (
    <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-5">
      <Tabs defaultValue="7days">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
          <TabsList className="bg-slate-800">
            <TabsTrigger
              value="7days"
              className="data-[state=active]:bg-purple-600"
            >
              Last 7 Days
            </TabsTrigger>
            <TabsTrigger
              value="30days"
              className="data-[state=active]:bg-purple-600"
            >
              Last 30 Days
            </TabsTrigger>
          </TabsList>

          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="text-slate-400">
                  <Info className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Average score across all players of this specialization</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>

        {/* Coming Soon content for all tabs */}
        <TabsContent value="7days" className="mt-0">
          <div className="h-64 w-full bg-slate-800/40 rounded-md border border-slate-700 flex items-center justify-center">
            <div className="text-center">
              <div className="text-slate-400 mb-2">Score Evolution Chart</div>
              <Badge
                variant="outline"
                className="text-purple-400 border-purple-600 text-lg py-2 px-4"
              >
                Coming Soon
              </Badge>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="30days" className="mt-0">
          <div className="h-64 w-full bg-slate-800/40 rounded-md border border-slate-700 flex items-center justify-center">
            <div className="text-center">
              <div className="text-slate-400 mb-2">
                Score Evolution Chart (30 Days)
              </div>
              <div className="text-slate-500">Coming Soon</div>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="season" className="mt-0">
          <div className="h-64 w-full bg-slate-800/40 rounded-md border border-slate-700 flex items-center justify-center">
            <div className="text-center">
              <div className="text-slate-400 mb-2">
                Score Evolution Chart (Season)
              </div>
              <div className="text-slate-500">Coming Soon</div>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default ScoreEvolution;
