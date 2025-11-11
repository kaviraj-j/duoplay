import {
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Typography,
  Box,
  Paper,
  Tooltip,
} from "@mui/material";
import { useState } from "react";
import { roomApi } from "@/api/room";
import constants from "@/contants";
import { ContentCopy } from "@mui/icons-material";
import { useRoom } from "@/contexts/RoomContext";

const ChooseOpponent = () => {
  const [open, setOpen] = useState(false);
  const [roomId, setRoomId] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [linkCopied, setLinkCopied] = useState(false);
  const { saveRoom } = useRoom();

  const handleOpen = () => setOpen(true);
  const handleClose = () => {
    setOpen(false);
    if (roomId) {
      // Save room to context
      roomApi.getRoom(roomId).then((room) => {
        saveRoom(room);
      });
    }
  };

  const generateRoomJoinLink = (roomId: string) =>
    `${constants.WEB_BASE_URL}/room/${roomId}/join`;

  const handleLinkCopy = async () => {
    try {
      await navigator.clipboard.writeText(generateRoomJoinLink(roomId ?? ""));
      setLinkCopied(true);
      setTimeout(() => setLinkCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy text:", err);
    }
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
    console.log("Play with random player - not implemented yet");
  };

  return (
    <>
      <Button
        variant="contained"
        onClick={handleOpen}
        className="!bg-blue-600 hover:!bg-blue-700 !text-white !font-semibold !px-6 !py-2.5 !rounded-lg !shadow-md"
      >
        Choose Opponent
      </Button>

      <Dialog
        open={open}
        onClose={handleClose}
        maxWidth="sm"
        fullWidth
        PaperProps={{
          className:
            "!rounded-2xl !p-4 !shadow-lg !border !border-gray-200 !bg-white",
        }}
      >
        <DialogTitle className="!text-xl !font-semibold !text-gray-800 !pb-0">
          Choose Your Opponent
        </DialogTitle>

        <DialogContent className="mt-3">
          {!roomId ? (
            <Box className="flex flex-col gap-3 mt-4">
              <Button
                variant="outlined"
                size="large"
                onClick={handlePlayWithFriend}
                disabled={isCreating}
                className="py-3 text-gray-800 font-medium border-gray-300 hover:bg-gray-50 rounded-lg"
              >
                {isCreating ? "Creating Room..." : "Play with a Friend"}
              </Button>

              <Button
                variant="outlined"
                size="large"
                onClick={handlePlayWithRandom}
                className="py-3 text-gray-800 font-medium border-gray-300 hover:bg-gray-50 rounded-lg"
              >
                Play with Random Player
              </Button>
            </Box>
          ) : (
            <Box className="text-center py-5">
              <Typography
                variant="h6"
                className="!text-green-700 !font-semibold !mb-2"
              >
                Room Created Successfully!
              </Typography>

              <Typography variant="body1" className="!text-gray-600 !mb-3">
                Share this joining link with your friend:
              </Typography>

              <Paper
                elevation={0}
                className="relative flex items-center justify-between bg-gray-50 border border-gray-300 rounded-lg p-3 text-sm text-gray-700"
              >
                <span className="truncate pr-10">
                  {generateRoomJoinLink(roomId)}
                </span>

                <Tooltip title={linkCopied ? "Copied!" : "Copy to clipboard"}>
                  <Button
                    onClick={handleLinkCopy}
                    className="!min-w-0 !p-1 hover:!bg-gray-100"
                  >
                    <ContentCopy fontSize="small" className="text-gray-600" />
                  </Button>
                </Tooltip>
              </Paper>

              {linkCopied && (
                <Typography
                  variant="body2"
                  className="!text-green-600 !mt-2 !font-medium"
                >
                  Link copied to clipboard!
                </Typography>
              )}
            </Box>
          )}
        </DialogContent>

        <DialogActions className="!pt-0">
          {roomId && (
            <Button
              onClick={handleClose}
              className="!bg-blue-600 hover:!bg-blue-700 !text-white !font-medium !px-4 !py-2 !rounded-md"
            >
              Close
            </Button>
          )}
        </DialogActions>
      </Dialog>
    </>
  );
};

export default ChooseOpponent;
