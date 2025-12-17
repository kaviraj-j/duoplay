import { gameApi, roomApi } from "@/api";
import { useRoom } from "@/contexts/RoomContext";
import type { GameListPayload } from "@/types";
import { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Container,
  Box,
} from "@mui/material";

const HomePage = () => {
  const [games, setGames] = useState<GameListPayload[]>([]);
  const { room } = useRoom();

  useEffect(() => {
    gameApi.getGamesList().then((gamesResponse) => {
      setGames(gamesResponse.data);
    });
  }, []);

  const handleChooseGame = (gameName: string) => {
    console.log("Choose game:", gameName);
    roomApi.chooseGame(room!.id, gameName);
  };

  return (
    <Container sx={{ py: 4 }}>
      <Box
        sx={{
          display: "grid",
          gridTemplateColumns: {
            xs: "1fr",
            sm: "repeat(2, 1fr)",
            md: "repeat(3, 1fr)",
          },
          gap: 3,
        }}
      >
        {games.map((game) => (
          <Card
            key={game.name}
            sx={{ height: "100%", display: "flex", flexDirection: "column" }}
          >
            <CardContent sx={{ flexGrow: 1 }}>
              <Typography variant="h5" component="h2" gutterBottom>
                {game.display_name}
              </Typography>
              {/* {game.description && (
                <Typography variant="body2" color="text.secondary">
                  {game.description}
                </Typography>
              )} */}
            </CardContent>
            <CardActions sx={{ justifyContent: "flex-end", p: 2 }}>
              <Button
                size="small"
                variant="outlined"
                onClick={() => handleChooseGame(game.name)}
                disabled={!room}
              >
                Choose Game
              </Button>
              {/* <Button
                size="small"
                variant="contained"
                onClick={() => handleViewGame(game.name)}
                disabled={!room}
              >
                View Game
              </Button> */}
            </CardActions>
          </Card>
        ))}
      </Box>
    </Container>
  );
};

export default HomePage;
