import { useState } from "react";
import { useAuthContext } from "@/contexts/AuthContext";
import {
  Button,
  Dialog,
  DialogTitle,
  DialogActions,
  IconButton,
} from "@mui/material";
import { Close as CloseIcon } from "@mui/icons-material";
import LogoutIcon from "@mui/icons-material/Logout";
import { useRoom } from "@/contexts/RoomContext";

export function LogoutModal() {
  const [isOpen, setIsOpen] = useState(false);
  const { logout } = useAuthContext();
  const { removeRoom } = useRoom();

  return (
    <>
      <Button
        variant="contained"
        color="primary"
        onClick={() => setIsOpen(true)}
      >
        <LogoutIcon />
      </Button>

      <Dialog
        open={isOpen}
        onClose={() => setIsOpen(false)}
        maxWidth="xs"
        fullWidth
      >
        <DialogTitle>
          Logout from DuoPlay?
          <IconButton
            aria-label="close"
            onClick={() => setIsOpen(false)}
            sx={{
              position: "absolute",
              right: 8,
              top: 8,
            }}
          >
            <CloseIcon />
          </IconButton>
        </DialogTitle>

        <DialogActions>
          <Button onClick={() => setIsOpen(false)}>Cancel</Button>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            onClick={() => {
              logout();
              removeRoom();
            }}
          >
            Confirm
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
