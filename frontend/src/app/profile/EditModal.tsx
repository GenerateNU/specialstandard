import { MultiSelect } from "@/components/ui/multiselect";
import { useSchools } from "@/hooks/useSchools";
import { X } from "lucide-react";
import React, { useEffect, useState } from "react";

interface EditModalProps {
  isOpen: boolean;
  onClose: () => void;
  therapistId: string;
  initialData: {
    first_name: string;
    last_name: string;
    school_names: string[];
    district_id?: number;
  };
  onSave: (
    therapistId: string,
    data: {
      first_name: string;
      last_name: string;
      school_ids: number[];
      district_id: number;
    }
  ) => void;
}

const EditModal: React.FC<EditModalProps> = ({
  isOpen,
  onClose,
  therapistId,
  initialData,
  onSave,
}) => {
  const { schools, districts } = useSchools();
  const [formData, setFormData] = useState({
    first_name: initialData.first_name,
    last_name: initialData.last_name,
  });
  const [selectedSchools, setSelectedSchools] = useState<string[]>([]);
  const [isInitialized, setIsInitialized] = useState(false);

  // Initialize schools from initialData
  useEffect(() => {
    if (isOpen && !isInitialized && schools.length > 0) {
      // Set schools
      if (initialData.school_names.length > 0) {
        const matchedSchoolIds = schools
          .filter((school) =>
            initialData.school_names.includes(school.name ?? "")
          )
          .map((school) => school.id!.toString());
        setSelectedSchools(matchedSchoolIds);
      }

      setIsInitialized(true);
    }

    // Reset when modal closes
    if (!isOpen && isInitialized) {
      setIsInitialized(false);
    }
  }, [isOpen, isInitialized, schools.length, initialData]);

  const handleSubmit = () => {
    if (!initialData.district_id || selectedSchools.length === 0) {
      return;
    }

    onSave(therapistId, {
      first_name: formData.first_name,
      last_name: formData.last_name,
      school_ids: selectedSchools.map((id) => Number(id)),
      district_id: initialData.district_id,
    });
    onClose();
  };

  // Filter schools based on the user's current district
  const filteredSchools = initialData.district_id
    ? schools.filter((school) => school.district_id === initialData.district_id)
    : schools;

  const schoolOptions = filteredSchools.map((school) => ({
    label: school.name ?? "",
    value: school.id!.toString(),
  }));

  // Get district name for display
  const districtName =
    districts.find((d) => d.id === initialData.district_id)?.name ??
    "No district assigned";

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 backdrop-blur-sm bg-black/30 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-2xl shadow-xl max-w-md w-full p-6">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-black">Edit Profile</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors cursor-pointer"
          >
            <X size={24} />
          </button>
        </div>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              First Name
            </label>
            <input
              type="text"
              value={formData.first_name}
              onChange={(e) =>
                setFormData({ ...formData, first_name: e.target.value })
              }
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-gray-500 focus:border-transparent text-gray-900"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Last Name
            </label>
            <input
              type="text"
              value={formData.last_name}
              onChange={(e) =>
                setFormData({ ...formData, last_name: e.target.value })
              }
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-gray-500 focus:border-transparent text-gray-900"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              District
            </label>
            <div className="w-full cursor-not-allowed px-4 py-2 bg-gray-100 border border-gray-300 rounded-lg text-gray-600">
              {districtName}
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Registered Schools
            </label>
            <MultiSelect
              options={schoolOptions}
              value={selectedSchools}
              onValueChange={setSelectedSchools}
              placeholder="Select Schools"
              showTags={false}
              showCount={false}
              className="w-full bg-white text-gray-600 border border-gray-300 rounded-lg focus:ring-2 focus:ring-gray-600 focus:border-transparent [&_button]:bg-white [&_button]:text-gray-600 [&_button]:border-gray-300 cursor-pointer"
            />
            {selectedSchools.length > 0 && (
              <p className="text-xs text-gray-500 mt-2">
                {selectedSchools.length} school
                {selectedSchools.length !== 1 ? "s" : ""} selected
              </p>
            )}
          </div>

          {selectedSchools.length > 0 && (
            <div className="mt-2">
              <button
                type="button"
                onClick={() => {
                  setSelectedSchools([]);
                }}
                className="text-xs text-gray-600 hover:text-gray-700 underline cursor-pointer"
              >
                Clear School Selection
              </button>
            </div>
          )}

          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="button"
              onClick={handleSubmit}
              disabled={selectedSchools.length === 0}
              className="flex-1 px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Save Changes
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default EditModal;
