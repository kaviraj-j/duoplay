import { useState } from "react";
import { useAuthContext } from "@/contexts/AuthContext";
import { userApi } from "@/api";
import {
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
} from "@mui/material";
import { Close as CloseIcon } from "@mui/icons-material";

export function LoginModal() {
  const [isOpen, setIsOpen] = useState(false);
  const [name, setName] = useState("");
  const { login } = useAuthContext();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const sanitizedName = name.trim();
      if (!sanitizedName) {
        return;
      }
      const res = await userApi.register({ name: sanitizedName });
      login(res.user, res.token);
      setIsOpen(false);
    } catch (error) {
      console.error("Login failed:", error);
    }
  };

  return (
    <>
      <Button
        variant="contained"
        color="primary"
        onClick={() => setIsOpen(true)}
      >
        Login
      </Button>

      <Dialog
        open={isOpen}
        onClose={() => setIsOpen(false)}
        maxWidth="xs"
        fullWidth
      >
        <DialogTitle>
          Login to DuoPlay
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

        <form onSubmit={handleSubmit}>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              label="Your Name"
              type="text"
              fullWidth
              variant="outlined"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </DialogContent>

          <DialogActions>
            <Button onClick={() => setIsOpen(false)}>Cancel</Button>
            <Button type="submit" variant="contained" color="primary">
              Join Duoplay
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </>
  );
}
