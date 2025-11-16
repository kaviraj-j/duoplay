import { useState } from "react";
import {
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Typography,
  Box,
} from "@mui/material";
import { Close as CloseIcon } from "@mui/icons-material";
import { roomApi } from "@/api";
import { useRoom, type GameChoiceData } from "@/contexts/RoomContext";

interface GameChoiceModalProps {
  open: boolean;
  gameChoice: GameChoiceData | null;
  onClose: () => void;
}

const gameDisplayNames: Record<string, string> = {
  tictactoe: "Tic Tac Toe",
  // Add more games as needed
};

export function GameChoiceModal({
  open,
  gameChoice,
  onClose,
}: GameChoiceModalProps) {
  const { room } = useRoom();
  const [isProcessing, setIsProcessing] = useState(false);

  const handleAccept = async () => {
    if (!gameChoice || !room) return;

    setIsProcessing(true);
    try {
      roomApi.acceptGame(room.id, gameChoice.game_type);
      onClose();
    } catch (error) {
      console.error("Failed to accept game:", error);
    } finally {
      setIsProcessing(false);
    }
  };

  const handleReject = async () => {
    if (!gameChoice || !room) return;

    setIsProcessing(true);
    try {
      roomApi.rejectGame(room.id, gameChoice.game_type);
      onClose();
    } catch (error) {
      console.error("Failed to reject game:", error);
    } finally {
      setIsProcessing(false);
    }
  };

  if (!gameChoice) return null;

  const gameDisplayName =
    gameDisplayNames[gameChoice.game_type] || gameChoice.game_type;

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="sm"
      fullWidth
      PaperProps={{
        className:
          "!rounded-2xl !p-4 !shadow-lg !border !border-gray-200 !bg-white",
      }}
    >
      <DialogTitle className="!text-xl !font-semibold !text-gray-800 !pb-2">
        Game Choice Request
        <IconButton
          aria-label="close"
          onClick={onClose}
          disabled={isProcessing}
          sx={{
            position: "absolute",
            right: 8,
            top: 8,
          }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>

      <DialogContent>
        <Box className="flex flex-col gap-4 mt-2">
          <Typography variant="body1" className="!text-gray-700">
            <strong>{gameChoice.player_name}</strong> has chosen to play{" "}
            <strong>{gameDisplayName}</strong>.
          </Typography>
          <Typography variant="body2" className="!text-gray-500">
            Would you like to accept or reject this game choice?
          </Typography>
        </Box>
      </DialogContent>

      <DialogActions className="!px-6 !pb-4 !pt-2">
        <Button
          onClick={handleReject}
          disabled={isProcessing}
          variant="outlined"
          color="error"
          className="!min-w-[100px]"
        >
          Reject
        </Button>
        <Button
          onClick={handleAccept}
          disabled={isProcessing}
          variant="contained"
          color="primary"
          className="!min-w-[100px]"
        >
          {isProcessing ? "Processing..." : "Accept"}
        </Button>
      </DialogActions>
    </Dialog>
  );
}

