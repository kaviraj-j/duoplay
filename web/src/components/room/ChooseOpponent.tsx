import { Button, Dialog, DialogTitle, DialogContent, DialogActions, Typography, Box } from "@mui/material";
import React, { useState } from "react";
import { roomApi } from "@/api/room";

const ChooseOpponent = () => {
  const [open, setOpen] = useState(false);
  const [roomId, setRoomId] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);

  const handleOpen = () => setOpen(true);
  const handleClose = () => {
    setOpen(false);
    setRoomId(null);
  };

  const handlePlayWithFriend = async () => {
    try {
      setIsCreating(true);
      const { roomId: newRoomId } = await roomApi.createRoom();
      setRoomId(newRoomId);
    } catch (error) {
      console.error("Failed to create room:", error);
    } finally {
      setIsCreating(false);
    }
  };

  const handlePlayWithRandom = () => {
    // TODO: Implement random player matching
    console.log("Play with random player - not implemented yet");
  };

  return (
    <>
      <Button variant="contained" onClick={handleOpen}>
        Choose Opponent
      </Button>

      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>Choose Your Opponent</DialogTitle>
        <DialogContent>
          {!roomId ? (
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
              <Button
                variant="outlined"
                size="large"
                onClick={handlePlayWithFriend}
                disabled={isCreating}
                sx={{ py: 2 }}
              >
                {isCreating ? "Creating Room..." : "Play with a Friend"}
              </Button>
              
              <Button
                variant="outlined"
                size="large"
                onClick={handlePlayWithRandom}
                sx={{ py: 2 }}
              >
                Play with Random Player
              </Button>
            </Box>
          ) : (
            <Box sx={{ textAlign: 'center', py: 3 }}>
              <Typography variant="h6" gutterBottom>
                Room Created Successfully!
              </Typography>
              <Typography variant="body1" sx={{ mb: 2 }}>
                Share this room ID with your friend:
              </Typography>
              <Typography 
                variant="h5" 
                sx={{ 
                  backgroundColor: '#f5f5f5', 
                  padding: 2, 
                  borderRadius: 1,
                  fontFamily: 'monospace',
                  fontWeight: 'bold'
                }}
              >
                {roomId}
              </Typography>
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          {roomId && (
            <Button onClick={handleClose} color="primary">
              Close
            </Button>
          )}
        </DialogActions>
      </Dialog>
    </>
  );
};

export default ChooseOpponent;
