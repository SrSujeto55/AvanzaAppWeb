import React from 'react';

import EditBaseUser from "../components/profiles/editBaseUser.tsx"
import EditBodilyMeasurements from "../components/profiles/EditBodilyMeasurements.tsx"
import ProtectedRoute from "../components/protectedRoute.tsx";
import { SidebarTrainer } from "../components/SideBar/SidebarTrainer.tsx"

export default function EditAthleteProfile() {

	return (
	<>
	<ProtectedRoute />

	<div className="flex flex-row">
		<SidebarTrainer />
		<div className="bg-blue-100 p-4 w-full">
			<EditBaseUser />
			<EditBodilyMeasurements />
		</div>
	</div>
	</>
	);
}
