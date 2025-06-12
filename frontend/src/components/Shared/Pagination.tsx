// Pagination.tsx - Composant réutilisable de pagination
import React from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange,
}) => {
  // Détermine quelles pages afficher
  const getPageNumbers = () => {
    const pageNumbers = [];
    const maxPagesToShow = 5;

    if (totalPages <= maxPagesToShow) {
      // Si nous avons moins de pages que maxPagesToShow, montrer toutes les pages
      for (let i = 0; i < totalPages; i++) {
        pageNumbers.push(i);
      }
    } else {
      // Sinon, montrer un sous-ensemble avec la page courante au milieu
      let startPage = Math.max(0, currentPage - Math.floor(maxPagesToShow / 2));
      let endPage = Math.min(totalPages - 1, startPage + maxPagesToShow - 1);

      // Ajuster si nous sommes près de la fin
      if (endPage - startPage < maxPagesToShow - 1) {
        startPage = Math.max(0, endPage - maxPagesToShow + 1);
      }

      for (let i = startPage; i <= endPage; i++) {
        pageNumbers.push(i);
      }

      // Ajouter les ellipses si nécessaire
      if (startPage > 0) {
        pageNumbers.unshift(-1); // Ellipsis at start
        pageNumbers.unshift(0); // First page
      }

      if (endPage < totalPages - 1) {
        pageNumbers.push(-2); // Ellipsis at end
        pageNumbers.push(totalPages - 1); // Last page
      }
    }

    return pageNumbers;
  };

  return (
    <div className="flex items-center justify-center mt-6 space-x-2">
      <button
        onClick={() => onPageChange(Math.max(0, currentPage - 1))}
        disabled={currentPage === 0}
        className="p-2 rounded-md bg-slate-800/50 text-white border border-slate-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-700"
        aria-label="Previous page"
      >
        <ChevronLeft size={16} />
      </button>

      {getPageNumbers().map((pageNum, index) => {
        if (pageNum === -1 || pageNum === -2) {
          // Ellipsis
          return (
            <span
              key={`ellipsis-${index}`}
              className="px-4 py-2 text-slate-400"
            >
              ...
            </span>
          );
        }

        return (
          <button
            key={pageNum}
            onClick={() => onPageChange(pageNum)}
            className={`px-4 py-2 rounded-md ${
              currentPage === pageNum
                ? "bg-purple-600 text-white"
                : "bg-slate-800/50 text-white border border-slate-700 hover:bg-slate-700"
            }`}
          >
            {pageNum + 1}
          </button>
        );
      })}

      <button
        onClick={() => onPageChange(Math.min(totalPages - 1, currentPage + 1))}
        disabled={currentPage === totalPages - 1}
        className="p-2 rounded-md bg-slate-800/50 text-white border border-slate-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-700"
        aria-label="Next page"
      >
        <ChevronRight size={16} />
      </button>
    </div>
  );
};

export default Pagination;
