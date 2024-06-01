import { CircleAlert } from "lucide-react";

export default function ErrorMessage({ err }: any) {
  return (
    <>
      {err?.message && (
        <div className="text-red-500 mb-3 text-sm flex items-center">
          <CircleAlert className="w-4 mr-1" />
          {err?.message}
        </div>
      )}
    </>
  );
}
