import { useState, useEffect } from 'react';
import { apiService } from '../services/api';

interface Team {
  id: string;
  name: string;
}

export function useTeam() {
  const [teams, setTeams] = useState<Team[]>([]);
  const [currentTeam, setCurrentTeam] = useState<Team | null>(null);

  useEffect(() => {
    fetchTeams();
  }, []);

  const fetchTeams = async () => {
    try {
      const fetchedTeams = await apiService.getTeams();
      setTeams(fetchedTeams);
      if (fetchedTeams.length > 0) {
        setCurrentTeam(fetchedTeams[0]);
      }
    } catch (error) {
      console.error('Failed to fetch teams:', error);
    }
  };

  const switchTeam = async (teamId: string) => {
    try {
      const team = await apiService.switchTeam(teamId);
      setCurrentTeam(team);
    } catch (error) {
      console.error('Failed to switch team:', error);
    }
  };

  return { teams, currentTeam, switchTeam };
}

