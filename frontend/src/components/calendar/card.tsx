import * as React from 'react';

//carta
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

import MyDrawer from './drawer.tsx';

export interface BasicCardProps {
    routineName: string;
    dayofWeek: number;
    workoutId: number;
    
  }

const bull = (
    <Box
      component="span"
      sx={{ display: 'inline-block', mx: '2px', transform: 'scale(0.8)' }}
    >
      •
    </Box>
  );
  
const BasicCard: React.FC<BasicCardProps> = ({ routineName, dayofWeek,workoutId }) => {
    return (
      <Card sx={{ minWidth: 275 }}>
        <CardContent>
          <Typography sx={{ fontSize: 14 }} color="text.secondary" gutterBottom>
            Rutina del día {dayofWeek}
          </Typography>
          <Typography variant="h5" component="div">
            {routineName}
          </Typography>
          <Typography sx={{ mb: 1.5 }} color="text.secondary">
            Duración: 1 hora
          </Typography>
          <Box sx={{ display: 'flex', justifyContent: 'center' }}>
            <MyDrawer idWorkout={workoutId}/>
          </Box>
        </CardContent>
      </Card>
    );
  };
  
  export default BasicCard;