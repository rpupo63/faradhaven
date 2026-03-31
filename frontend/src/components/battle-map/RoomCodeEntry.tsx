import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';

export const RoomCodeEntry: React.FC = () => {
  const [roomCode, setRoomCode] = useState('');
  const navigate = useNavigate();

  const handleJoin = () => {
    if (roomCode.trim()) {
      navigate(`/battle-map/room/${roomCode.trim()}`);
    }
  };

  return (
    <Card className="w-full max-w-md mx-auto mt-20 bg-card border-border text-foreground">
      <CardHeader>
        <CardTitle className="text-center text-2xl text-primary">Join Battle Map</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <Input
          placeholder="Enter Room Code"
          value={roomCode}
          onChange={(e) => setRoomCode(e.target.value)}
          className="bg-input border-input text-foreground"
        />
        <Button onClick={handleJoin} className="w-full bg-primary hover:bg-primary/90 text-primary-foreground">
          Join
        </Button>
      </CardContent>
    </Card>
  );
};
